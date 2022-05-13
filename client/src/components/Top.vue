<template>
  <div class="top">
    <img alt="Vue logo" src="../assets/logo.png">
    <h1>トップページ</h1>
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
.top {
  text-align: center;
  margin-top: 60px;
}
h1, h2 {
  font-weight: normal;
}
button {
  margin: 10px 0;
  padding: 10px;
}
</style>
