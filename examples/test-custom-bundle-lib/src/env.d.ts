declare module '.*?raw' {
  declare const raw: string
  export default raw
}

interface ImportMeta {
  env: {
    NODE_ENV: 'development' | 'production';
  };
}