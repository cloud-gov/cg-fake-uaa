#!/bin/sh

redirect_uri() {
    app_route=$2
    case $1 in 
        1) echo \'{\"redirect_uri\": [\"https://$app_route/auth/callback\"]}\'
        ;;
        2) echo \'{\"redirect_uri\": [\"https://$app_route/auth/callback\", \"https://$app_route/auth/logout\", \"https://$app_route/auth/callback\"]}\'
        #2) echo \'{\"redirect_uri\": [\"https://$app_route/auth/callback\", \"https://$app_route/auth/logout\"]}\'
        ;;
        *)
        echo "oops"
        exit 1
        ;;
    esac
}

up() {
    n=$1
    app="id-example-$n"
    cf push -m 128M --no-start --random-route $app
    app_route=$(cf apps | grep "^$app" | awk '{print $NF}')
    echo app_route: $app_route
    cf create-service cloud-gov-identity-provider oauth-client uaa-$app
    echo sleep 15 # it takes a moment to provision the oauth-client
    sleep 15 # it takes a moment to provision the oauth-client

    uris=`redirect_uri $n $app_route`
    set -x
    eval cf create-service-key uaa-$app uaa-$app-key -c "$uris"
    eval cf bind-service $app uaa-$app -c $uris
    set +x

    cf set-env $app UAA_AUTH_URL https://login.fr.cloud.gov/oauth/authorize
    cf set-env $app UAA_LOGOUT_URL https://login.fr.cloud.gov/logout.do
    cf set-env $app UAA_TOKEN_URL https://uaa.fr.cloud.gov/oauth/token
    cf start $app
}

down() {
    n=$1
    app="id-example-$n"
    cf delete -f $app
    cf delete-service-key -f uaa-$app uaa-$app-key
    cf delete-service -f uaa-$app
}

usage() {
    echo "$0 up|down 1|2"
    exit 1
}

case $2 in
  1|2) : ;;
  *) usage;;
esac

case $1 in 
  up) up $2;;
  down) down $2;;
  reset)
    down $2
    cf delete-orphaned-routes -f
    up $2;;
  *) usage;;
esac


# For LOGOUT:
# https://docs.cloudfoundry.org/api/uaa/version/4.19.0/index.html#logout-do
# curl 'http://localhost/logout.do?redirect....
# Note: If the chosen redirect URI is not whitelisted, users will land on the UAA login page. This is a security feature intended to prevent open redirects as per RFC 6749.



# But the authorize URL is:
# curl 'http://localhost/oauth/authorize?response...

    # you must provide -c option below or you get:
    #   Server error, status code: 502, error code: 10001, message: Service broker error: Must pass JSON configuration with field "redirect_uri"
      #-c '{"redirect_uri": ["https://'$app_route'/auth/callback"]}'
    #cf create-service-key my-uaa-client my-service-key -c '{"redirect_uri": ["https://'$app_route'/auth/callback", "https://'$app_route'/auth/logout"]}'
    
    # binding the service is required, or you get:
    #    /home/vcap/app/example-client.js:16
    #   client_id = vcap_services["cloud-gov-identity-provider"][0].credentials.client_id ;
    #   2018-08-02T18:10:55.57-0400 [APP/PROC/WEB/0] ERR TypeError: Cannot read property '0' of undefined 

    # Further you must provide the `-c` or you get:
    #  "description": "Service broker error: Must pass JSON configuration with field \"redirect_uri\"",
    #cf bind-service id-example my-uaa-client -c '{"redirect_uri": ["https://'$app_route'/auth/callback", "https://'$app_route'/auth/logout"]}'