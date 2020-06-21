<template>
  <v-app id="inspire">
    <v-navigation-drawer
      v-model="drawer"
      app
    >
      <v-list dense>
        <v-list-item link>
          <v-list-item-action>
            <v-icon>mdi-home</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Home</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar
      app
      color="indigo"
      dark
    >
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title><v-icon medium>lock_open</v-icon> Stage Lock Panel</v-toolbar-title>
    </v-app-bar>

    <v-main>
      <v-container fluid>
        <v-row
          align="center"
          justify="center"
        >
          <v-col>
            <v-simple-table>
              <template v-slot:default>
                <thead>
                  <tr class="text-center">
                    <td>Name</td>
                    <td>Lock status</td>
                    <td>Lock owner</td>
                    <td>Comment</td>
                    <td>Last Action</td>
                    <td>Actions</td>
                  </tr>
                </thead>
                <tbody>
                  <tr class="text-center" v-for="item in stages" :key="item.name">
                    <td>{{item.name}}</td>
                    <td>
                      <v-sheet :color="item.locked ? 'red': 'green'" rounded :elevation="3">{{item.locked ? 'Locked': 'Free'}}</v-sheet>
                    </td>
                    <td>{{item.locked_by || '-'}}</td>
                    <td>{{item.comment || '-'}}</td>
                    <td>{{moment(item.updated_at).format('HH:mm:ss DD.MM.YYYY')}}</td>
                    <td>
                      <div>
                        <v-btn @click.stop="openLockDialog(item.name)" class="ma-1 green darken-1" small :disabled="item.locked">Lock</v-btn>
                        <v-btn @click.stop="openUnLockDialog(item.name)" class="ma-1 red darken-1" small :disabled="!item.locked">Unlock</v-btn>
                        <v-btn @click.stop="openHistory(item.name)" class="ma-1 blue darken-1" small>History</v-btn>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </template>
            </v-simple-table>
          </v-col>
        </v-row>
        <lock-dialog
          :active="lockDialog"
          :stage="currentStage"
          @close="closeLockDialog"
          @notify="showNotification"
          @refresh="fetchStagesStatuses"
        >
        </lock-dialog>
        <un-lock-dialog
          :active="unLockDialog"
          :stage="currentStage"
          @close="closeUnLockDialog"
          @notify="showNotification"
          @refresh="fetchStagesStatuses"
        >
        </un-lock-dialog>
        <history
          :active="history"
          :stage="currentStage"
          @close="closeHistory"
        >
        </history>
        <v-snackbar
          right
          top
          v-model="notification"
          :color="notificationProps.color"
          timeout="3000"
        >
          {{notificationProps.text}}
          <template v-slot:action="{ attrs }">
            <v-btn
              dark
              text
              v-bind="attrs"
              @click="closeNotification"
            >
              Close
            </v-btn>
          </template>
        </v-snackbar>
      </v-container>
    </v-main>
    <v-footer
      color="indigo"
      app
    >
      <span class="white--text">&copy; 2020</span>
    </v-footer>
  </v-app>
</template>

<script>
import moment from 'moment'
import History from '@/components/History.vue'
import LockDialog from '@/components/LockDialog.vue'
import UnLockDialog from '@/components/UnlockDialog.vue'

export default {
  data: () => ({
    drawer: false,
    moment: moment,
    backUrl: process.env.VUE_APP_BACKEND_URL,
    stages: [],
    lockDialog: false,
    unLockDialog: false,
    notification: false,
    history: false,
    notificationProps: {
      color: null,
      text: null
    },
    currentStage: null
  }),
  components: {
    LockDialog,
    UnLockDialog,
    History
  },
  methods: {
    closeNotification () {
      this.notification = false
      this.notificationProps.color = null
      this.notificationProps.text = null
    },
    showNotification (text, color) {
      this.notification = true
      this.notificationProps.color = color
      this.notificationProps.text = text
    },
    async fetchStagesStatuses () {
      let response = await this.$http.get('/stages')
      this.stages = response.data.data
    },
    closeLockDialog () {
      this.lockDialog = false
      this.currentStage = null
    },
    closeUnLockDialog () {
      this.unLockDialog = false
      this.currentStage = null
    },
    openLockDialog (stage) {
      this.lockDialog = true
      this.currentStage = stage
    },
    openUnLockDialog (stage) {
      this.unLockDialog = true
      this.currentStage = stage
    },
    openHistory (stage) {
      this.currentStage = stage
      this.history = true
    },
    closeHistory () {
      this.history = false
      this.currentStage = null
    }
  },
  async created () {
    await this.fetchStagesStatuses()
  }
}
</script>
