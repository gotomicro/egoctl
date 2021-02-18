import { Settings as ProSettings } from '@ant-design/pro-layout';

type DefaultSettings = Partial<ProSettings> & {
  pwa: boolean;
};

const proSettings: DefaultSettings = {
  "navTheme": "dark",
  "primaryColor": "#1890ff",
  "layout": "side",
  "contentWidth": "Fluid",
  "fixedHeader": false,
  "fixSiderbar": true,
  "title": "EGO Low Code",
  "pwa": false,
  "iconfontUrl": "",
  "menu": {
    "locale": false,
  },
  "headerHeight": 48,
  "headerRender": false,
  "footerRender": false
};

export type { DefaultSettings };

export default proSettings;
