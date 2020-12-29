typedef struct {
    const char *output;
    const char *err;
} Response;

Response *getCommandOutput(const char *command);
void *DestroyResponse(Response response);
