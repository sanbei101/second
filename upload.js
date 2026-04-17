import ci from "miniprogram-ci";
import path from "path";
import packageJson from "./package.json" with { type: "json" };

(async () => {
  try {
    console.log("开始初始化上传任务...");
    const project = new ci.Project({
      appid: "wxcf90b73a0b3985b1",
      type: "miniProgram",
      projectPath: path.resolve("dist/build/mp-weixin"),
      privateKeyPath: path.resolve("upload-1.key"),
      ignores: ["node_modules/**/*"],
    });

    console.log("开始执行上传...");
    const uploadResult = await ci.upload({
      project,
      version: packageJson.version || "1.0.0",
      desc: "CI 自动构建上传",
      setting: {
        es6: true,
        minify: true,
        minifyJS: true,
        minifyWXML: true,
        minifyWXSS: true,
        autoPrefixWXSS: true,
      },
      onProgressUpdate: console.log,
    });

    console.log("上传成功:\n", uploadResult);
  } catch (err) {
    console.error("上传失败:", err);
    process.exit(1);
  }
})();
