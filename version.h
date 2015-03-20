#ifndef UNITYSCOPE_GO_VERSION_H
#define UNITYSCOPE_GO_VERSION_H

#include <unity/scopes/Version.h>

// check that we have a compatible version of lib-unityscopes installed
static_assert(UNITY_SCOPES_VERSION_MAJOR == 0, "Major version of Unity scopes API mismatch");
static_assert(UNITY_SCOPES_VERSION_MINOR >= 6, "Minor version of Unity scopes API mismatch");
static_assert(UNITY_SCOPES_VERSION_MICRO >= 10, "Micro version of Unity scopes API mismatch");

#endif
