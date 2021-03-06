#include <stdio.h>
#include <unistd.h> 
#include <stdlib.h>
#include <errno.h>
#include <string.h>

int main (int argc, char *argv[]) {
    char *final_path = "CHALLENGE";
    // have to do this to keep the suid
    setregid(getegid(), getegid());
    setreuid(geteuid(), geteuid());
    int arrayLength = (4 + argc - 1);
    char** arguments = (char**)malloc( sizeof(char*) * arrayLength);
    char* user = "--user=THE_USER";
    char* null = (char*)0;
    int i;
    arguments[0] = "sudo";
    arguments[1] = user;
    arguments[2] = final_path;
    for(i=3 ; i<arrayLength-1 ; ++i){
        arguments[i] = argv[i-2];
    }
    arguments[arrayLength-1] = null;
    int ret = execv("/usr/bin/sudo", arguments);
    free(arguments);
    if(ret != 0){
        fprintf(stderr, "%s\n", strerror(errno));
    }
    return ret;
}
