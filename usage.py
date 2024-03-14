import base64
import ctypes
import json
import os


def find_lib():
    dir_path = os.path.dirname(os.path.realpath(__file__))
    if os.path.isfile(os.path.join(dir_path, 'tlslib.dll')):
        filename = 'tlslib.dll'
    elif os.path.isfile(os.path.join(dir_path, 'tlslib.so')):
        filename = 'tlslib.so'
    elif os.path.isfile(os.path.join(dir_path, 'tlslib.dylib')):
        filename = 'tlslib.dylib'
    else:
        return None

    return os.path.join(dir_path, filename)


lib = ctypes.cdll.LoadLibrary(find_lib())

go_request = lib.libRequest
go_request.argtypes = [ctypes.c_char_p]
go_request.restype = ctypes.c_void_p

go_free = lib.libFree
go_free.argtypes = [ctypes.c_void_p]
go_free.restype = None


def request(url, method='GET', proxy='', body='', headers=None, timeout=10000, follow_redirects=True):
    if headers is None:
        headers = {}

    req = json.dumps({"url": url, "method": method, "proxy": proxy, "body": base64.b64encode(body.encode()).decode(),
                      "headers": headers, "timeout": timeout, "follow_redirects": follow_redirects})
    req_ptr = ctypes.create_string_buffer(req.encode('utf8'))
    resp_ptr = go_request(req_ptr)
    resp = json.loads(ctypes.cast(resp_ptr, ctypes.c_char_p).value)
    go_free(resp_ptr)

    if 'body' in resp and resp['body'] is not None:
        resp['body'] = base64.b64decode(resp['body'])

    return resp


def main():
    headers = {
        # The user agent must be valid or at least contain Chrome or Firefox
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36',
    }

    r = request("https://ja3er.com/json", headers=headers)
    print(r)


if __name__ == '__main__':
    main()
