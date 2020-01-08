#! /bin/bash

retrieve (){
    while [ 1 ]
    do
        echo "Grabbing Authentication token"
        token=$(cat /run/conjur/access-token)
        echo "Here's the access token:"
        echo "$token"
        echo "Formatting token"
        formatted_token=$(cat /run/conjur/access-token | base64 | tr -d '\r\n')
        echo "Formatted token is:"
        echo "$formatted_token"
        echo "Making call to retrieve $SECRET"
        output=$(curl -k -s -X GET -H "Authorization: Token token=\"$formatted_token\"" $CONJUR_APPLIANCE_URL/secrets/$CONJUR_ACCOUNT/variable/$SECRET)
        echo "Here's the output:"
        echo "$output"
        echo "-----"
        sleep 5
    done
}

retrieve
