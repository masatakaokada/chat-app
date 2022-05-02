<template>
  <div class="signin">
    <h2>Sign in</h2>
    <input
      v-model="email"
      type="text"
      placeholder="email"
    >
    <input
      v-model="password"
      type="password"
      placeholder="Password"
    >
    <button @click="signIn">
      Signin
    </button>
    <p>
      You don't have an account?
      <router-link to="/signup">
        create account now!!
      </router-link>
    </p>
  </div>
</template>

<script>
import { getAuth, signInWithEmailAndPassword } from 'firebase/auth'
export default {
  name: 'SignIn',
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
        this.$router.push('/chat')
      }, err => {
        alert(err.message)
      })
    }
  }
}
</script>

<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
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
