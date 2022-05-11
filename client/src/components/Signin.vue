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
    signIn: async function () {
      const auth = getAuth()
      try {
        const res = await signInWithEmailAndPassword(auth, this.email, this.password)
        localStorage.setItem('jwt', res.user.accessToken)
        this.$router.push('/')
      } catch (error) {
        switch (error.code) {
        case 'auth/wrong-password':
          alert('パスワードが違います')
          break
        case 'auth/invalid-email':
          alert('無効のメールアドレスです')
          break
        case 'auth/user-not-found':
          alert('ユーザーが存在しません')
          break
        default:
          alert(error.message)
          break
        }
      }
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
