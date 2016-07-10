"use strict";

const querystring = require("querystring");
const crypto = require("crypto");

const express = require("express");
const request = require("request");
const jwt = require("jsonwebtoken");

const PORT = process.env.PORT || 8000;
const CLIENT_ID = process.env.CLIENT_ID || 'my-client-id';
const CLIENT_SECRET = process.env.CLIENT_SECRET || 'my-client-secret';
const UAA_AUTH_URL = process.env.UAA_AUTH_URL || 
                     'http://localhost:8080/oauth/authorize';
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

// Exchanges an authorization code for an authorization token
function exchangeCodeForAuthToken(code, callback) {
  request.post(UAA_TOKEN_URL, {
    form: {
      code,
      grant_type: 'authorization_code',
      response_type: 'token',
      client_id: CLIENT_ID,
      client_secret: CLIENT_SECRET,
    },
  }, callback);
}

// Our "home page" contains a login link.
app.get('/', (req, res) => {
  const url = UAA_AUTH_URL + '?' + querystring.stringify({
    'client_id': CLIENT_ID,
    'state': FAKE_STATE,
    'response_type': 'code'
  });
  res.send(`<a href="${url}">Log in</a> via <code>${UAA_AUTH_URL}</code>!`);
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

  exchangeCodeForAuthToken(req.query.code, (err, response, body) => {
    if (err) {
      res.status(400).send(err);
      return;
    }

    const responseBody = JSON.parse(body);
    const decodedToken = jwt.decode(responseBody.access_token);
    const userEmail = decodedToken.email;

    res.send(`Hello, ${userEmail}! You have successfully authenticated.`);
  });
});

app.listen(PORT, () => {
  console.log("Listening on " + PORT);
});
