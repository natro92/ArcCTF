import request from "@/utils/request";

export function userRegisterService (registerData) {
    const params = new URLSearchParams();
    for (let key in registerData) {
        params.append(key, registerData[key]);
    }
    return request.post('/user/register', params);
}
