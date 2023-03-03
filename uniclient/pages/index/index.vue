<template>
	<view class="uni-padding-wrap uni-common-mt">
		<textarea id="share_url" placeholder-style="color:#ccc" :placeholder="placeholder"
			v-model="shareInfo"></textarea>
		<view style="display: flex;flex-direction: row; max-height: 40px; margin-top: 1%;">
			<!-- #ifdef !H5 -->
			<button type="default" plain="true" style="width: 50%;" @click="paste">粘贴</button>
			<!-- #endif  -->
			<button type="primary" style="width: 100%;margin-left: 1%;" @click="parseShareInfo" :loading="video.parsing"
				v-show="video.pss">解析</button>
			<button type="primary" style="width: 100%;margin-left: 1%;" @click="pickWenAn"
				v-show="video.wenan">提取文案</button>
			<button type="default" style="width: 80%;margin-left: 1%;" :disabled="video.download"
				@click="download">{{video.process}}</button>
		</view>
		<video crossOrigin="anonymous" id="myVideo" :src="video.url" @error="videoErrorCallback" enable-danmu danmu-btn
			controls></video>
	</view>
</template>

<script>
	import "@/static/clipboard.js"
	import * as comm from "@/static/req.js"
	export default {
		data() {
			return {
				title: 'Hello',
				shareInfo: "",
				placeholder: "粘贴视频分享地址",
				video: {
					title: "",
					parsing: false,
					url: "",
					cover: "",
					download: true,
					pss: true,
					wenan: false,
					process: "保存",
				},
			}
		},
		onLoad() {
			// uni.setStorageSync('storage_key', 'hello');
			return;
		},
		methods: {
			copy() {
				uni.setClipboardData({
					data: 'hello',
					success: () => {
						uni.showToast({
							icon: 'none',
							title: '复制成功'
						})
					},
					complete: () => {
						console.log("completed.");
					}
				});
			},
			paste() {
				console.log("paste click...")
				uni.getClipboardData({
					success: (res) => {
						this.shareInfo = res.data;
					},
					fail: (e) => {
						uni.showToast({
							icon: 'error',
							title: e
						})
					}
				})
			},
			parseShareInfo(e) {
				this.video.parsing = true;
				let that = this;
				comm.sendSignReqest(comm.baseHost + '/video/parse/share?url=' + encodeURIComponent(this.shareInfo), comm
					.Get, {},
					function(data) {
						that.video.parsing = false;
						if (data.code != 0) {
							uni.showToast({
								title: data.error,
								icon: "error",
							})
						} else {
							if (data.code == 0) {
								that.video.title = data.data.title;
								that.video.cover = data.data.cover_url;
								that.video.url = data.data.video_url;
								that.shareInfo = "";
								that.placeholder = "解析成功，可继续提取视频文案"
								that.video.download = false;
								that.video.pss = false;
								that.video.wenan = true;
								uni.showToast({
									icon: "success",
									title: "解析成功"
								})
							} else {
								uni.showToast({
									icon: "error",
									title: data.error
								})
							}
						}
					});
			},
			// 提取文案
			pickWenAn() {
				uni.showToast({
					title: "敬请期待",
					icon: "success",
				})
			},
			download() {
				//  #ifdef H5
				window.open(this.video.url);
				return
				// #endif
				const downloadTask = uni.downloadFile({
					url: this.video.url,
					header: {

					},
					success: (data) => {
						console.log(data);
						if (data.statusCode === 200) {
							uni.saveImageToPhotosAlbum({
								filePath: data.tempFilePath,
								success: function() {
									uni.showToast({
										title: "已保存相册",
										icon: "success"
									});
								},
								fail: function() {
									uni.showToast({
										title: "保存失败",
										icon: "error"
									});
								}
							});
						}
					}
				});

				downloadTask.onProgressUpdate((res) => {
					//console.log('已经下载的数据长度' + res.totalBytesWritten);
					//console.log('下载进度：' + res.progress);
					this.video.process = res.progress + "%";
					if (res.progress == 100) {
						this.video.process = "保存"
					}
				});
			},
			videoErrorCallback(e) {

			}
		}
	}
</script>

<style>
	@import url("../../common/uni.css");

	.content {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
	}

	.logo {
		height: 200rpx;
		width: 200rpx;
		margin-top: 200rpx;
		margin-left: auto;
		margin-right: auto;
		margin-bottom: 50rpx;
	}

	video {
		margin-top: 5%;
		width: 100%;
		object-fit: cover;
	}

	text-area {
		display: flex;
		min-width: 80%;
		min-height: 50%;
		border: #8f8f94 1rpx solid;
		background-color: aliceblue;
	}

	.title {
		font-size: 36rpx;
		color: #8f8f94;
	}
</style>
