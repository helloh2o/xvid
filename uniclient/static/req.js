import {
	hex_md5
} from "./md5";
// sign key
const key = "TpGLnL2VQHqV0pPu8ZLk3yBc5s@2023";
export const channel = "wx_min";
export const baseHost = "http://localhost:8080";
export const Get = "GET";
export const Post = "POST";
// 签名
function getSign(val) {
	return hex_md5(val);
}

// 发送带签名的请求
export function sendSignReqest(reqUrl, m, params, callback) {
	let h = {
		"_xt": "",
		"_ts": parseInt(new Date().getTime() / 1000),
		"_rs": randomString(16),
		"_xc": channel
	}
	h._xt = getSign(channel + key + h._rs + h._ts);
	console.log("sign:" + h._xt);
	uni.request({
		url: reqUrl,
		data: params,
		method: m,
		header: {
			'_xt': h._xt,
			'_ts': h._ts,
			'_rs': h._rs,
			'_xc': h._xc,
			'content-type': 'application/x-www-form-urlencoded',
		},
		success: (res) => {
			console.log(res.data);
			callback(res.data);
		},
		fail: (e) => {
			console.error(e);
		}
	});
}


function randomString(len) {
	len = len || 32;
	var $chars = 'ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz02345678oOLl9gqVvUuI1^><|';
	var maxPos = $chars.length;
	var ret = '';
	for (let i = 0; i < len; i++) {
		ret += $chars.charAt(Math.floor(Math.random() * maxPos));
	}
	return ret;
}
