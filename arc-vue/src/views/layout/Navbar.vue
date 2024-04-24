<script setup>
import {useStore, mapState, mapMutations} from "vuex";
import {computed} from "vue";
import {useRouter} from "vue-router";

const router = useRouter();
const store = useStore();

const userInfo = computed(() => store.state.userModule.userinfo);

function logout() {
    store.dispatch('userModule/logout');
}
</script>

<template>
    <nav>
        <div class="">
            <div class="flex justify-between h-16 px-10 shadow items-center">
                <div class="flex items-center space-x-8">
                    <h1 class="text-xl lg:text-2xl font-bold cursor-pointer" @click.prevent="router.replace({name: 'home'})">Arc</h1>
                    <div class="hidden md:flex justify-around space-x-4">
                        <div class="cursor-pointer hover:text-indigo-600 text-gray-700" @click.prevent="router.replace({name: 'home'})">HomePage</div>
                        <div class="cursor-pointer hover:text-indigo-600 text-gray-700" @click.prevent="router.replace({name: 'about'})">About</div>
                    </div>
                </div>
                <div>
                    <div v-if="userInfo" class="flex space-x-4 items-center">
                        <div> {{userInfo.name}} </div>
                        <div class="dropdown dropdown-end">
                            <div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar">
                                <div class="w-10 rounded-full ring ring-primary">
                                    <img alt="Tailwind CSS Navbar component" src="https://daisyui.com/images/stock/photo-1534528741775-53994a69daeb.jpg" />
                                </div>
                            </div>
                            <ul tabindex="0" class="mt-3 z-[1] p-2 shadow menu menu-sm dropdown-content bg-base-100 rounded-box w-52">
                                <li>
                                    <a class="justify-between" @click.prevent="router.push({name: 'profile'})">
                                        Profile
                                        <span class="badge">New</span>
                                    </a>
                                </li>
                                <li><a>Settings</a></li>
                                <li @click="logout"><a>Logout</a></li>
                            </ul>
                        </div>
                    </div>
                    <div v-if="!userInfo" class="flex space-x-4 items-center">
                        <div v-if="$route.name !== 'login'" @click.prevent="router.replace({name: 'login'})" class="cursor-pointer text-gray-800 text-sm">登录</div>
                        <div v-if="$route.name !== 'register' " @click.prevent="router.replace({name: 'register'})" class="cursor-pointer bg-indigo-600 px-4 py-2 rounded text-white hover:bg-indigo-500 text-sm">注册</div>
                    </div>
                </div>
            </div>
        </div>
    </nav>

</template>

<style scoped lang="scss">

</style>
