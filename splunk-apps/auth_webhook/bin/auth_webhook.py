import sys
import json
import csv
import gzip
import base64
from collections import OrderedDict
from future.moves.urllib.request import urlopen, Request
from future.moves.urllib.error import HTTPError, URLError

def send_webhook_request(url, body, basic_auth=None, bearer_token=None, user_agent=None):
    open = urlopen
    if url is None:
        sys.stderr.write("ERROR No URL provided\n")
        return False
    sys.stderr.write("INFO Sending POST request to url=%s with size=%d bytes payload\n" % (url, len(body)))
    sys.stderr.write("DEBUG Body: %s\n" % body)
    try:
        if sys.version_info >= (3, 0) and type(body) == str:
            body = body.encode()
        headers = {"Content-Type": "application/json", "User-Agent": user_agent}
        if basic_auth is not None:
            user, passwd = basic_auth
            headers['Authorization'] = b'Basic ' + base64.b64encode(f'{user}:{passwd}'.encode())
        if bearer_token is not None:
            headers['Authorization'] = "Bearer " + bearer_token
        req = Request(url, body, headers)
        res = open(req)
        if 200 <= res.code < 300:
            sys.stderr.write("INFO Webhook receiver responded with HTTP status=%d\n" % res.code)
            return True
        else:
            sys.stderr.write("ERROR Webhook receiver responded with HTTP status=%d\n" % res.code)
            return False
    except HTTPError as e:
        sys.stderr.write("ERROR Error sending webhook request: %s\n" % e)
    except URLError as e:
        sys.stderr.write("ERROR Error sending webhook request: %s\n" % e)
    except ValueError as e:
        sys.stderr.write("ERROR Invalid URL: %s\n" % e)
    return False


if __name__ == "__main__":
    if len(sys.argv) < 2 or sys.argv[1] != "--execute":
        sys.stderr.write("FATAL Unsupported execution mode (expected --execute flag)\n")
        sys.exit(1)
    try:
        settings = json.loads(sys.stdin.read())
        url = settings['configuration'].get('url')
        body = OrderedDict(
            sid=settings.get('sid'),
            search_name=settings.get('search_name'),
            app=settings.get('app'),
            owner=settings.get('owner'),
            results_link=settings.get('results_link'),
            result=settings.get('result')
        )
        user_agent = settings['configuration'].get('user_agent', 'Splunk')
        bearer_token = settings['configuration'].get('bearer_token', None)
        basic_auth = None
        basic_auth_user = settings['configuration'].get('basic_auth_user', None)
        basic_auth_pass = settings['configuration'].get('basic_auth_pass', None)
        if basic_auth_user is not None and basic_auth_pass is not None:
            basic_auth = (basic_auth_user, basic_auth_pass)
        if not send_webhook_request(url, json.dumps(body), basic_auth=basic_auth, bearer_token=bearer_token, user_agent=user_agent):
            sys.exit(2)
    except Exception as e:
        sys.stderr.write("ERROR Unexpected error: %s\n" % e)
        sys.exit(3)
