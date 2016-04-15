#define Py_LIMITED_API
#include <Python.h>

static PyMethodDef s3gitMethods[] = {  
    {NULL, NULL, 0, NULL}
};

static struct PyModuleDef s3gitmodule = {  
   PyModuleDef_HEAD_INIT, "s3git", NULL, -1, s3gitMethods
};

PyMODINIT_FUNC  
PyInit_s3git(void)  
{
    return PyModule_Create(&s3gitmodule);
}
