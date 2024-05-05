"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.globalConfig = void 0;
var globalConfig = {
  // 后端 api 根地址
  // 测试环境
  rootPath: 'http://localhost:8080',
  // 线上环境
  // rootPath : '',
  // 主页标题
  rootTitle: 'Surveillance-guy 监控小子',
  // 日志页的最大显示长度
  LogMaxLength: 20480
};
exports.globalConfig = globalConfig;