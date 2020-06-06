<template>
  <div id="app">
    <b-navbar class="is-dark">
      <template slot="brand">
            <b-navbar-item>
                <img
                    src="./assets/logo.png"
                    alt="Stage Lock Panel"
                >
                <p>Stage Lock Panel</p>
            </b-navbar-item>
        </template>
    </b-navbar>
    <main>
      <div class="container">
        <div class="section has-text-centered">
          <h1 class="title is-1">Stage servers status panel</h1>
        </div>
        <table class="table is-fullwidth">
          <thead>
            <td class="has-text-centered has-text-weight-semibold">Name</td>
            <td class="has-text-centered has-text-weight-semibold">Lock status</td>
            <td class="has-text-centered has-text-weight-semibold">Lock owner</td>
            <td class="has-text-centered has-text-weight-semibold">Comment</td>
            <td class="has-text-centered has-text-weight-semibold">Actions</td>
          </thead>
          <tbody>
            <tr v-for="item in stages" :key="item.name">
              <td class="has-text-centered">{{item.name}}</td>
              <td
                class="has-text-centered"
                :class="[item.locked ? 'has-background-danger': 'has-background-success']"
              >{{item.locked ? 'Locked': 'Free'}}</td>
              <td class="has-text-centered">{{item.locked_by || '-'}}</td>
              <td class="has-text-centered">{{item.comment || '-'}}</td>
              <td class="has-text-centered">
                <div class="buttons" style="justify-content: center;">
                  <b-button @click="lockDialog(item.name)" type="is-success" :disabled="item.locked">Lock</b-button>
                  <b-button @click="unlockDialog(item.name)" type="is-danger" :disabled="!item.locked">Unlock</b-button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </main>
    <footer class="footer has-background-dark has-text-white">
    </footer>
  </div>
</template>

<script>
import LockDialog from '@/components/LockDialog.vue'
import UnLockDialog from '@/components/UnlockDialog.vue'

export default {
  name: 'App',
  data () {
    return {
      backUrl: process.env.VUE_APP_BACKEND_URL,
      stages: []
    }
  },
  methods: {
    async lockDialog (stageName) {
      this.$buefy.modal.open({
        parent: this,
        component: LockDialog,
        hasModalCard: true,
        trapFocus: true,
        props: {
          stage: stageName
        },
        events: {
          refresh: this.fetchStagesStatuses
        }
      })
    },
    async unlockDialog (stageName) {
      this.$buefy.modal.open({
        parent: this,
        component: UnLockDialog,
        hasModalCard: true,
        trapFocus: true,
        props: {
          stage: stageName
        },
        events: {
          refresh: this.fetchStagesStatuses
        }
      })
    },
    async fetchStagesStatuses () {
      let response = await this.$http.get('/stages')
      this.stages = response.data.data
    }
  },
  async created () {
    await this.fetchStagesStatuses()
  }
}
</script>

<style lang="css">

  #app {
    display: flex;
    min-height: 100vh;
    flex-direction: column;
  }
  main {
    flex: 1;
  }

</style>
