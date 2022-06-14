const {defineConfig} = require('@vue/cli-service')
module.exports = defineConfig({
    transpileDependencies: true,
    publicPath: process.env.NODE_ENV === 'production' ? '/public' : '/',
    // outputDir: 'dist/public/',
    chainWebpack: config => {
        config
            .plugin('html')
            .tap(args => {
                args[0].title = 'gocron - 定时任务系统'
                return args
            })
    },
    devServer: {
        proxy: {
            '/api': {
                target: 'http://localhost:5920/api', //对应自己的接口
                changeOrigin: true,
                ws: true,
                pathRewrite: {
                    '^/api': ''
                }
            }
        }
    },
    configureWebpack:(config) => {
        if (process.env.NODE_ENV === 'production'){
            config.mode = 'production'
            config['performance'] = {//打包文件大小配置
                "maxEntrypointSize": 10000000,
                "maxAssetSize": 30000000
            }
        }
    }
})
