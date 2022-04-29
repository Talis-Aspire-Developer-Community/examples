from datetime import datetime
from email import parser
from pprint import pprint

import requests
from oauthlib.oauth2 import BackendApplicationClient, OAuth2Error
from requests_oauthlib import OAuth2Session
import requests as req
from logbook import Logger, FileHandler, StreamHandler
from argparse import ArgumentParser
import os
import sys
import json

log = Logger('MyFirstAPIScript')


# This script contains multiple classes to help keep the code clean.
# Normally you might have them in separate files and use Python's usual import mechanisms
# There are
class PersonaClient:
    """
    Handles the getting of an access token to make API calls with.
    """
    def __init__(self):
        self.token = ''
        # set these in your shell environment
        #  export ACTIVE_TALIS_PERSONA_ID='your-id'
        self.client_id = os.environ.get('ACTIVE_TALIS_PERSONA_ID')
        self.client_secret = os.environ.get('ACTIVE_TALIS_PERSONA_SECRET')
        # EU persona token url
        self.token_url = 'https://users.talis.com/oauth/tokens'
        # Canada should use the following
        # self.token_url = 'https://users.ca.talis.com/oauth/tokens'

        # this function is called when you instantiate the class and will get the token
        self.get_token_using_client_credentials()

    def get_token_using_client_credentials(self):
        """
        We use an OAuth client to get a token using client_credentials grant.
        :return:
        """
        if self.client_id and self.client_secret:
            try:
                # make a request to get a token
                log.info("Trying to get token for API access")

                client = BackendApplicationClient(client_id=self.client_id)
                oauth = OAuth2Session(client=client)
                response = oauth.fetch_token(self.token_url, client_id=self.client_id,
                                             client_secret=self.client_secret)
                self.token = response['access_token']
                # if we get and can set the token, then return True
                return True
            except OAuth2Error as e:
                log.error("Unable to get token. {} {}".format(e.status_code, e.description))
        else:
            log.error("Client ID or Secret are not set")
            return False

    def get_token(self):
        # In this example we are omitting handling of cases where
        # the token may have expired and you need to get a new token.
        # It makes more sense for you to use a token until it has expired
        # as this is less network traffic for your script
        # For this example we are just making one TARL API call later and so only going to call this once anyway!
        return self.token


class TARLAPIClient:
    """
    Provides an API client to interact with Talis Aspire APIs
    """
    def __init__(self, persona_client, tenant, talis_guid):
        self.persona_client = persona_client
        self.tenant = tenant
        self.talis_guid = talis_guid
        pass

    def set_authorization_header(self, headers):
        headers['Authorization'] = "Bearer {}".format(self.persona_client.get_token())
        return headers

    def get_default_headers(self):
        """
        These headers will be added to every request
        :return:
        """
        default_headers = {
            'Accept': 'application/json',
            'X-Effective-User': self.talis_guid
        }
        return self.set_authorization_header(default_headers)

    def get_list(self, list_id):
        """
        Implements the get list route. https://rl.talis.com/3/docs#operation/getList
        https://rl.talis.com/3/{shortCode}/lists/{listGuid}

        :param list_id:
        :return:
        """
        url = "https://rl.talis.com/3/{}/lists/{}".format(self.tenant, list_id)
        # Canada should comment out the line above and use this line below.
        # url = "https://rl.ca.talis.com/3/{}/lists/{}".format(self.tenant, list_id)
        try:
            headers = self.get_default_headers()
            res = req.get(url, headers=headers)
        except Exception as e:
            log.error("Could not make request: {}".format(e))

        return json.loads(res.content)


# Main processing using the above classes and their methods
def __main__():
    # set up some logging
    # register some logging handlers (uses https://logbook.readthedocs.io/en/stable/)
    log_handler = FileHandler(
        "get_list_title-{}.log".format(datetime.strftime(datetime.now(), '%Y%m%d%H%M%S')),
        mode='w',
        level='INFO',
        bubble=True
    )
    stdout_handler = StreamHandler(sys.stdout, level='INFO', bubble=True)

    # setup some arguments that we'll pass on the command line
    parser = ArgumentParser()
    parser.add_argument("-t", "--tenant", dest="tenant",
                        help="The reading lists tenant short code")
    parser.add_argument("-l", "--list", dest="list_id",
                        help="The list id to lookup")

    args = parser.parse_args()
    
    print('This is an example script to get the title of a list:')
    # Any code run within these logging handler with statements will log messages to both your terminal and a file.
    with stdout_handler.applicationbound():
        with log_handler.applicationbound():
            persona_client = PersonaClient()
            tarl_client = TARLAPIClient(persona_client, args.tenant, os.environ.get('ACTIVE_TALIS_PERSONA_ID'))
            my_list = tarl_client.get_list(args.list_id)  # You'd pass in your own list ID here!
            
            log.debug("{}".format(my_list))

            log.info("Title: {}".format(my_list['data']['attributes']['title']))
            log.info("Last Published: {}".format(my_list['data']['attributes']['last_published']))

__main__()
