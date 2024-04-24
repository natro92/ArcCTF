import { createStore } from 'vuex';
import userModule from "@/store/module/user";

export default createStore({
  strict: process.env.NODE_ENV !== 'production',
  state: {
  },
  getters: {
  },
  mutations: {
  },
  actions: {
  },
  modules: {
    userModule,
  },
});
