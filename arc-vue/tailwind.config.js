/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{html,js,vue}'],
  theme: {
    extend: {},
  },
  // eslint-disable-next-line global-require,import/no-extraneous-dependencies
  plugins: [require('daisyui')],
  daisyui: {
    themes: [
      {
        themes: ['light'],
      },
    ],
  },
};
