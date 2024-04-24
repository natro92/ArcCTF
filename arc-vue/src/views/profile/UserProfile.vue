<script setup>
import {computed} from "vue";
import {useStore} from "vuex";

const store = useStore();


const userInfo = computed(() => store.state.userModule.userinfo);

</script>

<template>
    <div v-if="userInfo" class="flex space-x-4 items-center">
        <div>您好 {{userInfo.name}} </div>
    </div>
    <div>sdfsdf</div>
    <div>sdfsdf</div>
    <div>sdfsdf</div>
    <div>sdfsdf</div>
    <div>sdfsdf</div>
    <div>sdfsdf</div>
    <div>sdfsdf</div>
    <div>sdfsdf</div>
    <div>sdfsdf</div>
</template>

<style scoped lang="scss">

</style>
