
Register("zt220.new-prospect:before", (params, target) => {
    if (params["Name"]) {
        return Promise.resolve();
    }
    return new Promise(async (resolve, reject) => {
        params["Name"] = await Shell.inputbox("File name:");
        if (!params["Name"]) resolve(false);
        resolve(params);
    });
});

Register("zt220.new-app:before", (params, target) => {
    if (params["Name"]) {
        return Promise.resolve();
    }
    return new Promise(async (resolve, reject) => {
        params["Name"] = await Shell.inputbox("File name:");
        if (!params["Name"]) resolve(false);
        resolve(params);
    });
});

Register("zt220.new-page:before", (params, target) => {
    if (params["Name"]) {
        return Promise.resolve();
    }
    return new Promise(async (resolve, reject) => {
        params["Name"] = await Shell.inputbox("File name:");
        if (!params["Name"]) resolve(false);
        resolve(params);
    });
});