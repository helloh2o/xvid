<template>
	<view>
		<view class="uni-padding-wrap uni-common-mt">
			<uni-section :title="username" type="line">
				<uni-list>
					<uni-list-item title="剩余次数" :rightText="leftwt" />
					<!--
							<uni-list-item title="提取文案" showArrow link="redirectTo" to="./chat" @click="onClick" />
							<uni-list-item title="打开错误页面(查看控制台)" showArrow link="redirectTo" to="./chats" @click="onClick" />
							-->
				</uni-list>
			</uni-section>
		</view>
		<view class="version">
			<text class="is-text-box">当前版本 v0.0.1</text>
		</view>
	</view>
</template>

<script>
	import * as req from "@/static/req.js"
	const USER_INFO = "user_info";
	export default {
		data() {
			return {
				username: "",
				leftwt: 0,
			}
		},
		onLoad() {
			let userInfo = uni.getStorageSync(USER_INFO);
			if (!userInfo) {
				// #ifdef !H5
				uni.login({
					provider: 'weixin', //使用微信登录
					success: function(loginRes) {
						console.log(loginRes);
						let authCode = loginRes.code;
						if (authCode) {
							req.sendSignReqest(req.baseHost + "/user/wx_login", req.Post, {
								"code": authCode,
								"channel": req.channel
							}, function(data) {
								if (data.code != 0) {
									uni.showToast({
										title: data.error,
										icon: "error",
									})
								} else {
									uni.setStorageSync(USER_INFO, data);
								}
							});
						}
					}
				});
				// #endif
			} else {
				console.log("local:", userInfo)
				this.username = userInfo.data.nickname;
				this.leftwt = userInfo.data.left_times_wt;
			}
		},
		methods: {

		}
	}
</script>

<style>
	.version {
		text-align: center;
		left: 0;
		right: 0;
		position: fixed;
		bottom: 15rpx;
		font-weight: 200;
		color: grey;
	}
</style>
