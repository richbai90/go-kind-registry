// preauthorize.h
#pragma once
#include <Security/Security.h>


void LaunchPreauthorizedProcess(AuthorizationRef *authRef, void (*cb) ());
OSStatus PreauthorizePrivilegedProcess(AuthorizationRef *authRef);