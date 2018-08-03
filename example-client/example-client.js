"use strict";

const querystring = require("querystring");
const crypto = require("crypto");

const express = require("express");
const request = require("request");
const jwt = require("jsonwebtoken");

const PORT = process.env.PORT || 8000;
var client_id = 'my_client_id';
var client_secret = 'my_client_secret';
var root_url = 'https://locahost:' + PORT + '/'
var vcap_services = JSON.parse(process.env.VCAP_SERVICES || null);
var vcap_application = JSON.parse(process.env.VCAP_APPLICATION || null);

if ( vcap_services ) {
  client_id = vcap_services["cloud-gov-identity-provider"][0].credentials.client_id ;
  client_secret = vcap_services["cloud-gov-identity-provider"][0].credentials.client_secret ;
}

if (vcap_application) {
  root_url =  'https://' + vcap_application["application_uris"][0] + '/'
}

const CLIENT_ID = process.env.CLIENT_ID || client_id;
const CLIENT_SECRET = process.env.CLIENT_SECRET || client_secret;
const UAA_AUTH_URL = process.env.UAA_AUTH_URL || 
                     'http://localhost:8080/oauth/authorize';
const UAA_LOGOUT_URL = process.env.UAA_LOGOUT_URL || 
                     'http://localhost:8080/logout.do';
const UAA_TOKEN_URL = process.env.UAA_TOKEN_URL ||
                      'http://localhost:8080/oauth/token';

// A "real" implementation would create a cryptographically secure
// random string per session, but we're just a demo app so we'll use
// a random string tied to the process lifetime.
//
// For more information on why this is important, see:
//
// http://www.twobotechnologies.com/blog/2014/02/importance-of-state-in-oauth2.html
const FAKE_STATE = crypto.randomBytes(64).toString('hex');

const app = express();

// A real implementation would use session middleware; because we're
// an example app, we'll just use a global object.
let session = {};

// This posts to the token endpoint, either to exchange an authorization
// code for an access token, or to obtain a new access token from a
// refresh token (depending on the contents of the given payload).
//
// It then sets the session data accordingly: if the token request failed,
// the user (if any) is logged out. Otherwise, the user is logged in.
function postToTokenUrlAndSetSession(payload, callback) {
  request.post(UAA_TOKEN_URL, {
    form: payload,
  }, function(err, response, body) {
    if (!err && response && response.statusCode !== 200) {
      err = new Error(`got HTTP ${response.statusCode}`);
    }
    if (err) {
      session = {};
    } else {
      const responseBody = JSON.parse(body);
      const decodedToken = jwt.decode(responseBody.access_token);

      Object.assign(session, {
        email: decodedToken.email,
        refreshToken: responseBody.refresh_token,
        expiry: Date.now() + responseBody.expires_in * 1000
      });
    }
    callback(err);
  });
}

// This middleware detects if the user's access token has expired; if
// so, it attempts to renew it with the server. This is primarily a
// security measure, to ensure that users who lose access to the system
// are logged out of it ASAP.
app.use(function tokenRefreshMiddleware(req, res, next) {
  if (session.expiry && Date.now() > session.expiry) {
    // Our session has expired! Let's renew it.
    console.log('User session has expired, attempting to renew it.');
    postToTokenUrlAndSetSession({
      grant_type: 'refresh_token',
      refresh_token: session.refreshToken,
      client_id: CLIENT_ID,
      client_secret: CLIENT_SECRET,
    }, (err) => {
      if (err) {
        console.log('Renewal unsuccessful, user was logged out.');
      } else {
        console.log('Renewal successful.');
      }
      next();
    });
  } else {
    next();
  }
});

// HTML for our home page when the user is logged in.
function getLoggedInHtml() {
  const expiresIn = Math.floor((session.expiry - Date.now()) / 1000);
  return `
    <p>Hello ${session.email}!</p>
    <p>Your access token lasts for another ${expiresIn} seconds,
    but will be renewed automatically.</p>
    <p>You can also <a href="/auth/logout">logout</a>.</p>
  `;
}

// HTML for our home page when the user is logged out.
function getLoggedOutHtml() {
  const url = UAA_AUTH_URL + '?' + querystring.stringify({
    'client_id': CLIENT_ID,
    'state': FAKE_STATE,
    'response_type': 'code'
  });
  return `<a href="${url}">Log in</a> via <code>${UAA_AUTH_URL}</code>!`;
}

// Our home page.
app.get('/', (req, res) => {
  res.send(session.email ? getLoggedInHtml() : getLoggedOutHtml());
});

// Handler for the registered callback URL.
app.get('/auth/callback', (req, res) => {
  if (!req.query.code) {
    res.status(400).send('Missing "code" query parameter');
    return;
  }

  if (req.query.state !== FAKE_STATE) {
    res.status(400).send('Invalid "state"');
    return;
  }

  // Exchange our authorization code for an authorization token.
  postToTokenUrlAndSetSession({
    code: req.query.code,
    grant_type: 'authorization_code',
    response_type: 'token',
    client_id: CLIENT_ID,
    client_secret: CLIENT_SECRET,
  }, (err) => {
    if (err) {
      res.status(400).send(err);
      return;
    }

    res.redirect('/');
  });
});

// Simple logout view to clear our session.
app.get('/auth/logout', (req, res) => {
  if (session.email) { // user is authenticated
    const logout_url = UAA_LOGOUT_URL + '?' + querystring.stringify({
      'client_id': CLIENT_ID,
      'redirect': root_url
    });
    session = {};
    res.redirect(logout_url);
  } else {
    res.redirect('/');
  }
});

app.listen(PORT, () => {
  console.log("Listening on " + PORT);
});
