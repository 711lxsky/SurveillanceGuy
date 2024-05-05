// 这段可选代码用于注册一个服务工作者。
// register() 不会默认调用。

// 这使应用在后续访问时加载速度更快，并提供了离线功能。但是，这也意味着开发者（和用户）
// 只会在关闭当前页面上所有标签页后，在后续访问中看到部署的更新，因为之前缓存的资源在后台更新。
// 要了解更多关于这种模型的好处以及如何启用，请参阅 https://bit.ly/CRA-PWA

const localhostRegex = /^127(?:\.(?:25[0-5]|2[0-4]\d|1\d{2}|[1-9]\d|\d)){3}$/;

const isLocalhost = Boolean(
  window.location.hostname === 'localhost' ||
  window.location.hostname === '[::1]' ||
  localhostRegex.exec(window.location.hostname)
);

export async function register(config) {
  if (process.env.NODE_ENV === 'production' && 'serviceWorker' in navigator) {
    // URL 构造函数在所有支持 SW 的浏览器中可用。
    const publicUrl = new URL(process.env.PUBLIC_URL, window.location.href);
    if (publicUrl.origin !== window.location.origin) {
      // 如果 PUBLIC_URL 与页面服务的源不同，我们的服务工作者将无法正常工作。
      // 这可能发生在使用 CDN 服务静态资产时；参见 https://github.com/facebook/create-react-app/issues/2374
      return;
    }

    window.addEventListener('load', async () => {
      const swUrl = `${process.env.PUBLIC_URL}/service-worker.js`;

      if (isLocalhost) {
        // 运行在本地主机。让我们检查服务工作者是否仍然存在。
        await checkValidServiceWorker(swUrl, config);
      } else {
        // 不是本地主机。直接注册服务工作者
        await loadAndRegisterSW(swUrl, config);
      }
    });
  }
}

async function loadAndRegisterSW(swUrl, config) {
  try {
    const response = await fetch(swUrl);
    if (response.status === 404 || (!response.ok && response.status !== 0)) {
      throw new Error('服务工作者未找到');
    }
    const contentType = response.headers.get('content-type');
    if (!contentType || contentType.indexOf('javascript') === -1) {
      throw new Error('服务工作者不是 JS 文件');
    }
    await registerValidSW(swUrl, config);
  } catch (error) {
    console.error('服务工作者注册时出错:', error);
    if (window.location.reload) {
      window.location.reload();
    }
  }
}

async function registerValidSW(swUrl, config) {
  try {
    const registration = await navigator.serviceWorker.register(swUrl);
    registration.onupdatefound = () => {
      const installingWorker = registration.installing;
      if (installingWorker == null) {
        return;
      }
      installingWorker.onstatechange = () => {
        if (installingWorker.state === 'installed') {
          handleServiceWorkerStateChange(registration, config);
        }
      };
    };
  } catch (error) {
    console.error('服务工作者注册时出错:', error);
  }
}

function handleServiceWorkerStateChange(registration, config) {
  if (navigator.serviceWorker.controller) {
    console.log(
      '新内容可用，将在关闭所有页面标签页后使用。参阅 https://bit.ly/CRA-PWA。'
    );

    if (config ?. config.onUpdate) {
      config.onUpdate(registration);
    }
  } else {
    console.log('内容已缓存供离线使用。');

    if (config ?. config.onSuccess) {
      config.onSuccess(registration);
    }
  }
}

export function unregister() {
  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.ready.then(registration => {
      registration.unregister();
    });
  }
}

function checkValidServiceWorker(swUrl, config) {
  // 如果找不到服务工作者，此操作将刷新页面
  loadAndRegisterSW(swUrl, config);
}