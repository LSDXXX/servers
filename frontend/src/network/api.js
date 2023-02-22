import request from "./public";

//登录
export function authCodeLogin(params) {
  return request({
    url: "/login",
    method: "post",
    data: params,
  });
}
//退出
// export function authLogout() {
//   return request({
//     url: baseUrl + "/logout",
//     method: "get",
//   });
// }
//获取用户数据
// export function getUserInfo(params) {
//   return request({
//     url: baseUrl + "/getUserInfo",
//     method: "get",
//     params: qs.stringfy(params),
//   });
// }
