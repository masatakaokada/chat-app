<template>
  <div class="top">
    <h1>{{ msg }}</h1>
    <button @click="signOut">
      ログアウト
    </button>
    <button @click="$router.push('/chat')">
      全体チャット
    </button>
  </div>
</template>

<script>
import { getAuth, signOut } from 'firebase/auth'
export default {
  name: 'Top',
  data () {
    return {
      msg: 'トップページ'
    }
  },
  methods: {
    signOut: function () {
      const auth = getAuth();
      signOut(auth).then(() => {
        localStorage.removeItem('jwt')
        this.$router.push('/signin')
      })
    },
  }
}
</script>

<style scoped>
h1, h2 {
  font-weight: normal;
}
button {
  margin: 10px 0;
  padding: 10px;
}
</style>
