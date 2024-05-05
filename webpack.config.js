module.exports = {
    // ...
    module: {
      rules: [
        {
          test: /\.js$/,
          use: ["source-map-loader"],
          enforce: "pre",
          exclude: /node_modules\/(?!mutationobserver-shim)/, // 排除除 mutationobserver-shim 外的所有 node_modules
        },
      ],
    },
    // ...
  };
  