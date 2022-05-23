<template>
  <div class="fill-height room-new">
    <v-container class="fill-height d-flex justify-center align-center">
      <v-card flat min-width="300px" :width="mdAndUp ? '500px' : null">
        <v-card-title class="d-flex justify-center pa-8">
          <h4>ルーム作成</h4>
        </v-card-title>
        <v-divider />
        <div class="px-6 py-8">
          <v-form>
            <v-text-field v-model="name" label="ルーム名(20文字まで)" required hide-details />
            <div class="py-4">
              <div v-for="user in users" :key="user.id">
                <v-checkbox
                  v-model="checkedUserIds"
                  :label="user.name"
                  :value="user.id"
                  hide-details
                />
              </div>
            </div>
          </v-form>
          <div class="pb-8">
            <v-btn color="#FFCB00" height="48px" block @click="createRoom">
              作成する
            </v-btn>
          </div>
          <v-divider />
          <p class="pt-8">
            <v-btn height="48px" block @click="$router.push('/')">
              トップページ
            </v-btn>
          </p>
        </div>
      </v-card>
    </v-container>
  </div>
</template>

<script>
import axios from 'axios'
import { useDisplay } from 'vuetify'
export default {
  name: 'RoomNew',
  setup () {
    const { mdAndUp } = useDisplay()
    return { mdAndUp }
  },
  data () {
    return {
      users: [],
      name: '',
      checkedUserIds: []
    }
  },
  async created() {
    const res = await axios.get(`${process.env.VUE_APP_API_URL}/users`, {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('jwt')}` }
    })
    this.users = res.data
  },
  methods: {
    createRoom: async function() {
      try {
        await axios.post(`${process.env.VUE_APP_API_URL}/rooms`, { name: this.name, userIds: this.checkedUserIds }, {
          headers: { 'Authorization': `Bearer ${localStorage.getItem('jwt')}` }
        })
        this.$router.push('/')
      } catch (error) {
          alert(error.response.data.ValidationErrors.join('\n'))
      }
    }
  },
}
</script>

<style scoped>
.room-new {
  background-color: #EEEEEE
}
</style>
