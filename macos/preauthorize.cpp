// preauthorize.cpp

#include "preauthorize.h"
OSStatus PreauthorizePrivilegedProcess(AuthorizationRef *authRef)
{
    AuthorizationItem item = {kAuthorizationRightExecute, 0, NULL, 0};
    AuthorizationRights rights = {1, &item};
    AuthorizationFlags flags = kAuthorizationFlagInteractionAllowed | kAuthorizationFlagExtendRights | kAuthorizationFlagPreAuthorize;
    return AuthorizationCreate(&rights, kAuthorizationEmptyEnvironment, flags, authRef);
}

void LaunchPreauthorizedProcess(AuthorizationRef *authRef, void (*cb) ()) {
    cb();
    AuthorizationFree(*authRef, kAuthorizationFlagDestroyRights);
}