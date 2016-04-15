#
# Copyright 2016 Frank Wessels <fwessels@xs4all.nl>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

from io import BytesIO
import json
import shutil
import tempfile
from cffi import FFI
from ctypes import *
ffi = FFI()

ffi.cdef("""
typedef signed char GoInt8;
typedef unsigned char GoUint8;
typedef short GoInt16;
typedef unsigned short GoUint16;
typedef int GoInt32;
typedef unsigned int GoUint32;
typedef long long GoInt64;
typedef unsigned long long GoUint64;
typedef GoInt64 GoInt;

extern GoInt s3git_init_repository(char* p0);

extern GoInt s3git_open_repository(char* p0);

extern GoInt s3git_clone(char* p0, char* p1, char* p2, char* p3);

extern char* s3git_add(char* p0, char* p1);

extern GoInt s3git_commit(char* p0, char* p1);

extern char* s3git_get(char* p0, char* p1);

extern GoInt s3git_push(char* p0, GoUint8 p1);

extern GoInt s3git_pull(char* p0);

extern char* s3git_list(char* p0, char* p1);

extern char* s3git_list_commits(char* p0);
""")

__s3gitlib__ = ffi.dlopen("s3git-py.so")

class Repository(object):
    path = ""

    def __init__(self, path):
        self.path = path

    def add(self, contents):
        stream = contents
        if isinstance(contents, str):  # Convert to stream for strings
            stream = BytesIO(contents.encode('utf-8'))

        with tempfile.NamedTemporaryFile() as temp:
            shutil.copyfileobj (stream, temp)
            temp.flush()
            ret = __s3gitlib__.s3git_add(ffi.new("char[]", self.path.encode('utf-8')),
                                         ffi.new("char[]", temp.name.encode('utf-8')))
            return ffi.string(ret)

    def commit(self, message):
        ret = __s3gitlib__.s3git_commit(ffi.new("char[]", self.path.encode('utf-8')),
                                        ffi.new("char[]", message.encode('utf-8')))
        return ret

    def get(self, hash):
        ret = __s3gitlib__.s3git_get(ffi.new("char[]", self.path.encode('utf-8')),
                                     ffi.new("char[]", hash.encode('utf-8')))
        return open(ffi.string(ret), "rb")

    def push(self, hydrate=False):
        arg_hydrate = 1 if hydrate else 0
        ret = __s3gitlib__.s3git_push(ffi.new("char[]", self.path.encode('utf-8')),
                                      ffi.new("unsigned char", arg_hydrate))
        return ret

    def pull(self):
        ret = __s3gitlib__.s3git_pull(ffi.new("char[]", self.path.encode('utf-8')))
        return ret

    def list(self, hash):
        ret = __s3gitlib__.s3git_list(ffi.new("char[]", self.path.encode('utf-8')),
                                      ffi.new("char[]", hash.encode('utf-8')))
        return ffi.string(ret).decode('utf-8').split(',')

    def list_commits(self):
        json_str = __s3gitlib__.s3git_list_commits(ffi.new("char[]", self.path.encode('utf-8')))
        return json.loads(ffi.string(json_str).decode('utf-8'))

def init_repository(path):
    repo = Repository(path)
    __s3gitlib__.s3git_init_repository(ffi.new("char[]", repo.path.encode('utf-8')))
    return repo

def open_repository(path):
    repo = Repository(path)
    __s3gitlib__.s3git_open_repository(ffi.new("char[]", repo.path.encode('utf-8')))
    return repo

def clone(url, path, access="", secret=""):
    repo = Repository(path)
    arg_url    = ffi.new("char[]", url.encode('utf-8'))
    arg_path   = ffi.new("char[]", repo.path.encode('utf-8'))
    arg_access = ffi.new("char[]", access.encode('utf-8'))
    arg_secret = ffi.new("char[]", secret.encode('utf-8'))
    __s3gitlib__.s3git_clone(arg_url, arg_path, arg_access, arg_secret)
    return repo
