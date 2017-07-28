var webpack = require('webpack');
var path = require('path');
var HtmlWebpackPlugin = require('html-webpack-plugin');

var TARGET = process.env.TARGET;
var ROOT_PATH = path.resolve(__dirname);

config = {
    entry: [
        path.join(ROOT_PATH, 'src', 'main.jsx'),
    ],
    resolve: {
        extensions: ['*', '.js', '.jsx'],
    },
    output: {
        path: path.join(ROOT_PATH, 'build'),
        filename: 'bundle.js',
        publicPath: '/'
    },
    module: {
        loaders: [{
                test: /\.jsx$/,
                loader: 'babel-loader',
                exclude: /node_modules/,
                query: {
                    presets: ['es2015', 'react']
                },
                include: path.join(ROOT_PATH, 'src'),
            },
            {
                test: /favicon\.ico/,
                loader: "file?name=favicon.ico",
            },
            {
                test: /index\.ico/,
                loader: "file?name=favicon.ico",
            },
            {
                test: /\.css$/,
                loaders: ['style', 'css'],
            },
            {
                test: /\.woff2?(\?v=\d+\.\d+\.\d+)?$/,
                loader: "url?limit=10000&mimetype=application/font-woff"
            },
            {
                test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/,
                loader: "url?limit=10000&mimetype=application/octet-stream"
            },
            {
                test: /\.eot(\?v=\d+\.\d+\.\d+)?$/,
                loader: "file"
            },
            {
                test: /\.svg(\?v=\d+\.\d+\.\d+)?$/,
                loader: "url?limit=10000&mimetype=image/svg+xml",
            }
        ],
    },
    plugins: [
        new webpack.ProvidePlugin({
            $: "jquery",
            jQuery: "jquery",
            "window.jQuery": "jquery",
        }),
        new HtmlWebpackPlugin({
            template: 'src/index.html',
            inject: 'body',
            hash: true,
        }),
    ],
};

if (TARGET == "dev") {
    config.plugins.push(new webpack.DefinePlugin({
        '__BASE_URL__': JSON.stringify('http://localhost:8000'),
    }));
} else if (TARGET == "devhot") {
    config.entry.splice(0, 0, "webpack/hot/only-dev-server");
    //config.module.loaders[0].loaders.splice(0, 0, 'react-hot');
    config.plugins.push(new webpack.DefinePlugin({
        '__BASE_URL__': JSON.stringify('http://localhost:8000/'),
    }));
} else {
    config.plugins.push(new webpack.DefinePlugin({
        'process.env': {
            // This has effect on the react lib size
            'NODE_ENV': JSON.stringify('production'),
        },
        '__BASE_URL__': JSON.stringify(''),
    }));
    config.plugins.push(new webpack.optimize.UglifyJsPlugin({
        compress: {
            warnings: false
        },
    }));
}

module.exports = config;