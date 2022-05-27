<template>
  <div class="fill-height signup">
    <v-container class="fill-height d-flex justify-center align-center">
      <v-card flat>
        <v-card-title class="d-flex justify-center pa-8">
          <h4>新規登録</h4>
        </v-card-title>
        <v-divider />
        <div class="px-6 py-8">
          <v-form>
            <v-text-field v-model="name" label="ニックネーム" :rules="[rules.required, rules.counter]" counter maxlength="20" />
            <v-text-field v-model="email" label="メールアドレス" :rules="[rules.required]" />
            <v-text-field v-model="password" label="パスワード" :rules="[rules.required]" type="password" />
          </v-form>
          <div class="pb-8">
            <v-btn color="#FFCB00" height="48px" block @click="signUp">
              登録する
            </v-btn>
          </div>
          <v-divider />
          <p class="pt-8 px-8">
            既にアカウントをお持ちの方
            <router-link to="/signin">
              ログイン
            </router-link>
          </p>
        </div>
      </v-card>
    </v-container>
  </div>
</template>

<script>
import axios from 'axios'
import { getAuth, createUserWithEmailAndPassword, deleteUser } from 'firebase/auth'
export default {
  name: 'Signup',
  data () {
    return {
      name: '',
      email: '',
      password: '',
      rules: {
        required: value => !!value || '入力必須です',
        counter: value => value.length <= 20 || '20文字以内で入力してください',
      },
    }
  },
  methods: {
    signUp: async function () {
      const auth = getAuth();
      try {
        const res = await createUserWithEmailAndPassword(auth, this.email, this.password)
        try {
          await axios.post(`${process.env.VUE_APP_API_URL}/users`, { name: this.name }, {
            headers: { 'Authorization': `Bearer ${res.user.accessToken}` }
          })
        } catch (error) {
          await deleteUser(res.user)
          alert(error.message)
        }
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
        case 'auth/weak-password':
          alert('6文字以上でパスワードを設定してください')
          break
        case 'auth/email-already-in-use':
          alert('すでに存在しているメールアドレスです')
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
.signup {
  background-color: #EEEEEE
}
</style>
