#!/bin/sh

up() {
    cf push -m 128M --no-start --random-route id-example
    app_route=$(cf apps | grep '^id-example' | awk '{print $NF}')
    echo app_route: $app_route
    cf create-service cloud-gov-identity-provider oauth-client my-uaa-client

    echo sleep 15 # it takes a moment to provision the oauth-client
    sleep 15 # it takes a moment to provision the oauth-client

    # you must provide -c option below or you get:
    #   Server error, status code: 502, error code: 10001, message: Service broker error: Must pass JSON configuration with field "redirect_uri"
      #-c '{"redirect_uri": ["https://'$app_route'/auth/callback"]}'
    set -x
    cf create-service-key my-uaa-client my-service-key -c '{"redirect_uri": ["https://'$app_route'/auth/callback"]}'
    
    # binding the service is required, or you get:
    #    /home/vcap/app/example-client.js:16
    #   client_id = vcap_services["cloud-gov-identity-provider"][0].credentials.client_id ;
    #   2018-08-02T18:10:55.57-0400 [APP/PROC/WEB/0] ERR TypeError: Cannot read property '0' of undefined 

    # Further you must provide the `-c` or you get:
    #  "description": "Service broker error: Must pass JSON configuration with field \"redirect_uri\"",
    cf bind-service id-example my-uaa-client -c '{"redirect_uri": ["https://'$app_route'/auth/callback"]}'
    set +x

    cf set-env id-example UAA_AUTH_URL https://login.fr.cloud.gov/oauth/authorize
    cf set-env id-example UAA_LOGOUT_URL https://login.fr.cloud.gov/logout.do
    cf set-env id-example UAA_TOKEN_URL https://uaa.fr.cloud.gov/oauth/token
    cf start id-example
}

down() {
    cf delete -f id-example
    cf delete-service-key -f my-uaa-client my-service-key
    cf delete-service -f my-uaa-client
}

usage() {
    echo "$0 up|down"
    exit 1
}

case $1 in 
  up) up;;
  down) down;;
  reset)
    down
    up;;
  *) usage;;
esac


# For LOGOUT:
# https://docs.cloudfoundry.org/api/uaa/version/4.19.0/index.html#logout-do
# curl 'http://localhost/logout.do?redirect....
# Note: If the chosen redirect URI is not whitelisted, users will land on the UAA login page. This is a security feature intended to prevent open redirects as per RFC 6749.



# But the authorize URL is:
# curl 'http://localhost/oauth/authorize?response...