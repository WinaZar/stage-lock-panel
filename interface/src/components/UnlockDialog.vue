<template>
  <v-row justify="center">
    <v-dialog v-model="active" persistent max-width="600px">
      <v-card>
        <v-card-title>
          <span class="headline">Unlock <span class="font-weight-bold">{{stage}}</span></span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  label="Lock Code"
                  v-model.trim="$v.lockCode.$model"
                  :error-messages="fieldErrors.lockCode"
                >
                </v-text-field>
              </v-col>
            </v-row>
          </v-container>
          <small>*indicates required field</small>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="close">Close</v-btn>
          <v-btn color="blue darken-1" text @click="unlock">Unlock</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-row>
</template>

<script>
import { minLength, maxLength } from 'vuelidate/lib/validators'
export default {
  name: 'UnLockDialog',
  props: {
    stage: String,
    active: Boolean
  },
  data () {
    return {
      lockCode: ''
    }
  },
  validations: {
    lockCode: {
      minLength: minLength(3),
      maxLength: maxLength(8)
    }
  },
  computed: {
    fieldErrors () {
      let errors = {}
      for (let key of Object.keys(this.$v.$params)) {
        errors[key] = this.getErrorMessage(key)
      }
      return errors
    }
  },
  methods: {
    getErrorMessage (name) {
      let validation = this.$v[name]
      if (validation.$error) {
        let message = 'Invalid field value'
        if (validation.required === false) {
          message = 'Field is required'
        } else if (validation.minLength === false) {
          message = `Field must have at least ${validation.$params.minLength.min} letters`
        } else if (validation.maxLength === false) {
          message = `Field must have less than ${validation.$params.maxLength.max} letters`
        }
        return message
      }
    },
    close () {
      this.lockCode = ''
      this.lockOwner = ''
      this.lockComment = ''
      this.$emit('close')
    },
    async unlock () {
      this.$v.$touch()
      if (this.$v.$invalid) {
        this.$emit('notify', 'Check form fields', 'red')
        return false
      }
      let data = null
      try {
        let response = await this.$http.post(`/stages/${this.stage}/unlock`, {
          code: this.lockCode,
          locked_by: 'dummy'
        })
        data = response.data
      } catch (e) {
        data = e.response.data
      }
      this.$emit('notify', data.message, data.status === 'error' ? 'red' : 'green')
      this.$emit('refresh')
      this.close()
    }
  }
}
</script>
