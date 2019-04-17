# Testing ks-kcode AWS Lambda endpoint using Python requests
#
# Equivalent using curl:
# url -X POST -d @1022_pumpkins.kcode https://ygfqhyxu69.execute-api.us-west-1.amazonaws.com/Prod/?method=blocks --header "Content-Type:application/json"
# {"result":["events_onGesture","objects_scale","objects_scale","events_onGesture","objects_scale"]}
# curl -X POST -d @1022_pumpkins.kcode https://ygfqhyxu69.execute-api.us-west-1.amazonaws.com/Prod/?method=spells --header "Content-Type:application/json"
# {"result":["engorgio","reducio"]}
# curl -X POST -d @1022_pumpkins.kcode https://ygfqhyxu69.execute-api.us-west-1.amazonaws.com/Prod/?method=validate --header "Content-Type:application/json"
# {"result":["expectedSpells=2","foundSpells=2","expectedBlocks=5","foundBlocks=5","isValid=true"]}
#

import json
import requests

ROOT = 'https://ygfqhyxu69.execute-api.us-west-1.amazonaws.com/Prod/'

# === ::getRequest() ===
def getRequest(root,command,headers={},params={},verbose=False):
    # Make a GET request
    print(f"::getRequest(command={command},params={params})")
    #auth = requests.auth.HTTPBasicAuth(username,password)
    url = f"{root}{command}"
    custom_headers = {'Accept': 'application/json','Content-Type': 'application/json'}
    custom_headers = {**custom_headers,**headers}
    req = requests.Request('GET',url,params=params,headers=custom_headers)
    prepared = req.prepare()
    isPost = False
    #dumpRequest(prepared,isPost,verbose)
    s = requests.Session()
    r = s.send(prepared)
    return r

# === ::postRequest() ===
def postRequest(root,command,body,headers={},params={},verbose=False):
    # Make a POST request
    print(f"::postRequest(command={command},params={params},body={len(body)} bytes)")
    #auth = requests.auth.HTTPBasicAuth(username,password)
    url = f"{root}{command}"
    custom_headers = {'Accept': 'application/json','Content-Type': 'application/json'}
    custom_headers = {**custom_headers,**headers}
    req = requests.Request('POST',url,params=params,headers=custom_headers,data=body)
    prepared = req.prepare()
    isPost = True
    #dumpRequest(prepared,isPost,verbose)
    s = requests.Session()
    r = s.send(prepared)
    return r

if __name__ == '__main__':
    kcode = ''
    command = ""
    headers = {}
    params = {}
    verbose = False

    # 1. ------- Test GET request ------- 
    r = getRequest(ROOT,command,headers,params,verbose) 
    expected = {'result':'ok'}
    received = r.json()
    #print(json.dumps(received))
    #print(f"GET request - expecting '{expected}' (type={type(expected)}) and received '{received}' (type={type(received)})")
    assert(expected == received)

    # 2. ------- Test POST request: BLOCKS ------- 
    fname = '1022_pumpkins.kcode'
    with open(fname,'r') as f:
        kcode = f.read()
    body = kcode
    params = {'method':'blocks'}
    r = postRequest(ROOT,command,body,headers,params,verbose)
    expected = {'result': ['events_onGesture', 'objects_scale', 'objects_scale', 'events_onGesture', 'objects_scale']}
    received = r.json()
    assert(expected == received)
    fname = '020_big_beans.kcode'
    with open(fname,'r') as f:
        kcode = f.read()
    body = kcode
    params = {'method':'blocks'}
    r = postRequest(ROOT,command,body,headers,params,verbose)
    expected = {'result': ['events_whileFlick', 'objects_scale', 'wand_vibrate', 'speaker#speaker_play', 'speaker#speaker_sample', 'events_whileFlick', 'objects_scale', 'wand_vibrate', 'speaker#speaker_play', 'speaker#speaker_sample']}
    received = r.json()
    assert(expected == received)
    print(received)

    # 3. ------- Test POST request: SPELLS ------- 
    fname = '1022_pumpkins.kcode'
    with open(fname,'r') as f:
        kcode = f.read()
    body = kcode
    params = {'method':'spells'}
    r = postRequest(ROOT,command,body,headers,params,verbose)
    expected = {'result': ['engorgio', 'reducio']}
    received = r.json()
    assert(expected == received)

    # 4. ------- Test POST request: PARTS ------- 
    fname = '020_big_beans.kcode'
    with open(fname,'r') as f:
        kcode = f.read()
    body = kcode
    params = {'method':'parts'}
    r = postRequest(ROOT,command,body,headers,params,verbose)
    expected = {'result': ['speaker']}
    received = r.json()
    assert(expected == received)

    # 5. ------- Test POST request: SCENE ------- 
    fname = '1022_pumpkins.kcode'
    with open(fname,'r') as f:
        kcode = f.read()
    body = kcode
    params = {'method':'scene'}
    r = postRequest(ROOT,command,body,headers,params,verbose)
    expected = {'result': ['puzzle022']}
    received = r.json()
    assert(expected == received)

    # 6. ------- Test POST request: VALIDATE ------- 
    fname = '1022_pumpkins.kcode'
    with open(fname,'r') as f:
        kcode = f.read()
    body = kcode
    params = {'method':'validate'}
    r = postRequest(ROOT,command,body,headers,params,verbose)
    # blocks, spells, parts, scene len
    expected = {'result': ['5', '5', '2', '2', '0', '0', '0', '9','true']}
    received = r.json()
    print(received)
    assert(expected == received)
    fname = '020_big_beans.kcode'
    with open(fname,'r') as f:
        kcode = f.read()
    body = kcode
    params = {'method':'validate'}
    r = postRequest(ROOT,command,body,headers,params,verbose)
    # blocks, spells, parts, scene len
    expected = {'result': ['10', '10', '0', '0', '1', '1', '0', '15','true']}
    received = r.json()
    print(received)
    assert(expected == received)

    # 7. ------- Test POST request: invalid method ------- 
    fname = '1022_pumpkins.kcode'
    with open(fname,'r') as f:
        kcode = f.read()
    body = kcode
    params = {'method':'foo'}
    r = postRequest(ROOT,command,body,headers,params,verbose)
    expected = {'result': 'Method Not Allowed - allowed methods are: spells, blocks or validate'}
    received = r.json()
    assert(expected == received)
    
    print("-------- PASSED --------")