<template>
  <div class="signup">
    <h2>新規登録</h2>
    <input
      v-model="name"
      type="text"
      placeholder="ニックネーム"
    >
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
    <button @click="signUp">
      登録する
    </button>
    <p>
      既にアカウントをお持ちの方
      <router-link to="/signin">
        ログイン
      </router-link>
    </p>
  </div>
</template>

<script>
import axios from 'axios'
import { getAuth, createUserWithEmailAndPassword } from 'firebase/auth'
export default {
  name: 'Signup',
  data () {
    return {
      name: '',
      email: '',
      password: ''
    }
  },
  methods: {
    signUp: async function () {
      const auth = getAuth();
      try {
        const res = await createUserWithEmailAndPassword(auth, this.email, this.password)
        axios.post(`${process.env.VUE_APP_API_URL}/users`, { name: this.name }, {
          headers: { 'Authorization': `Bearer ${res.user.accessToken}` }
        })
        localStorage.setItem('jwt', res.user.accessToken)
        this.$router.push('/')
      } catch (error) {
        alert(error.message)
      }
    }
  }
}
</script>

<style scoped>
h2 {
  font-weight: normal;
}
.signup {
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
