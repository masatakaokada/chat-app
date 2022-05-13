module.exports = {
  env: {
    node: true,
  },
  extends: ["eslint:recommended", "plugin:vue/recommended"],
  rules: {
    "vue/multi-word-component-names": "off",
    "vue/max-attributes-per-line": "off",
  },
};
