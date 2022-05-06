<template>
  <div class="signin">
    <h2>ログイン</h2>
    <input
      v-model="email"
      type="text"
      placeholder="メールアドレス"
    >
    <input
      v-model="password"
      type="password"
      placeholder="パスワード"
    >
    <button @click="signIn">
      ログイン
    </button>
    <p>
      初めての登録ですか?
      <router-link to="/signup">
        新規登録
      </router-link>
    </p>
  </div>
</template>

<script>
import { getAuth, signInWithEmailAndPassword } from 'firebase/auth'
export default {
  name: 'Signin',
  data: function () {
    return {
      email: '',
      password: ''
    }
  },
  methods: {
    signIn: function () {
      const auth = getAuth()
      signInWithEmailAndPassword(auth, this.email, this.password).then(res => {
        localStorage.setItem('jwt', res.user.accessToken)
        this.$router.push('/')
      }, err => {
        alert(err.message)
      })
    }
  }
}
</script>

<style scoped>
h2 {
  font-weight: normal;
}
.signin {
  margin-top: 20px;
  display: flex;
  flex-flow: column nowrap;
  justify-content: center;
  align-items: center
}
input {
  margin: 10px 0;
  padding: 10px;
}
button {
  margin: 10px 0;
  padding: 10px;
}
</style>
