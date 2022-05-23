<template>
  <div class="fill-height signin">
    <v-container class="fill-height d-flex justify-center align-center">
      <v-card flat>
        <v-card-title class="d-flex justify-center pa-8">
          <h4>ログイン</h4>
        </v-card-title>
        <v-divider />
        <div class="px-6 py-8">
          <v-form>
            <v-text-field v-model="email" label="メールアドレス" required />
            <v-text-field v-model="password" label="パスワード" required type="password" />
          </v-form>
          <div class="pb-8">
            <v-btn color="#FFCB00" height="48px" block @click="signIn">
              ログイン
            </v-btn>
          </div>
          <v-divider />
          <p class="pt-8 px-8">
            初めての登録ですか?
            <router-link to="/signup">
              新規登録
            </router-link>
          </p>
        </div>
      </v-card>
    </v-container>
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
.signin {
  background-color: #EEEEEE
}
</style>
