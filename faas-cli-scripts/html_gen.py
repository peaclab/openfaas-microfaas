import sys
import random
import json

def handler(req):
    """handle a request to the function
    Args:
        req (str): request body
    """
    params = json.loads(req)
   
    n = int(params["n"])
    output = "<html><table><tr><th>Char</th><th>ASCII</th></tr>"
    for _ in range(n):
        ascii_code = random.getrandbits(8)
        output += "<tr><td>" + str(chr(ascii_code)) + "</td><td>" + str(ascii_code) + "</td></tr>"

    output += "</table></html>"

    print(output)
handler(sys.argv[1])