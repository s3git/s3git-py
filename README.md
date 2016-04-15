s3git-py
========

[![Join the chat at https://gitter.im/s3git/s3git](https://badges.gitter.im/s3git/s3git.svg)](https://gitter.im/s3git/s3git?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

This is a Python module for [s3git](https://github.com/s3git/s3git). Please see [here](https://github.com/s3git/s3git/blob/master/README.md) for a general introduction to s3git including its use cases.

This module is based on the [s3git-go](https://github.com/s3git/s3git-go) package that is invoked via foreign function interface (FFI).

**DISCLAIMER: This software is still under development (although the storage format/model is stable) -- use at your own peril for now**

**Note that the API is not stable yet, you can expect minor changes/extensions**

Installation
------------

Please make sure you have a working Golang environment installed, otherwise you cannot build the shared library. See [install golang](https://github.com/minio/minio/blob/master/INSTALLGO.md) for setting up a working Golang environment.

Also the [s3git-go](https://github.com/s3git/s3git-go) package needs to be available locally (`go get -d github.com/s3git/s3git-go`).

```sh
$ go build -buildmode=c-shared -o s3git-py.so
```

Examples
--------

Here are some common examples of typical tasks.

Create a repository
-------------------

```
$ python3
>>> import s3git
>>> repo = s3git.init_repository('.')
>>> repo.add('hello s3git')
>>> repo.commit('My first s3git commit from Python')
>>> exit()
$ s3git log --pretty
113b9e900588fd830deed131c0f2f2cdd88b4751286ecc4967d0d4990aff6821c2426425955ec104797b142e93204a604e4ba420b3b617724237476128333de8 My first s3git commit from Python
$ s3git ls
c518dc5f1d95258dc91f6d285e7ea7300f37dea4dd517173f2e23afe0cb52bc9d8eb18683cdcf377e96a2d5a81585e61f6d27fa5d017cad53836bd050e9f105f
$ s3git cat c518dc5f1d95258dc91f6d285e7ea7300f37dea4dd517173f2e23afe0cb52bc9d8eb18683cdcf377e96a2d5a81585e61f6d27fa5d017cad53836bd050e9f105f
hello s3git
```

Clone a repository
------------------

```py
>>> import s3git
>>> repo = s3git.clone('s3://s3git-spoon-knife', '.', access='AKIAJYNT4FCBFWDQPERQ', secret='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')
>>> repo.list('')
```

Dump contents of a repository
-----------------------------

```py
>>> import s3git
>>> repo = s3git.init_repository('.')
>>> repo.add('hello s3git')
>>> repo.add('Python rocks')
>>> repo.add(open('local-file.txt', 'rb'))
>>> for l in repo.list(''):
>>>     repo.get(l).read()
```

List multiple commits
---------------------

```py
>>> import s3git
>>> repo = s3git.init_repository('.')
>>> repo.add('string 1')
>>> repo.commit('first commit')
>>> repo.add('string 2')
>>> repo.commit('second commit')
>>> for c in repo.list_commits():
>>>     print(c['Message'])
```

Limitations and Optimizations
-----------------------------

- To be described 

Contributions
-------------

Contributions are welcome, please submit a pull request for any enhancements.

License
-------

s3git-rb is released under the Apache License v2.0. You can find the complete text in the file LICENSE.
