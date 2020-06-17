<template>
  <v-row justify="center">
    <v-dialog v-model="active" persistent max-width="600px">
      <v-card>
        <v-card-title>
          <span class="headline">Lock <span class="font-weight-bold">{{stage}}</span></span>
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
              <v-col cols="12">
                <v-text-field
                  label="Lock Owner*"
                  v-model.trim="$v.lockOwner.$model"
                  :error-messages="fieldErrors.lockOwner"
                  required
                >
                </v-text-field>
              </v-col>
              <v-col cols="12">
                <v-textarea
                  label="Lock Comment"
                  v-model.trim="$v.lockComment.$model"
                  :error-messages="fieldErrors.lockComment"
                >
                </v-textarea>
              </v-col>
            </v-row>
          </v-container>
          <small>*indicates required field</small>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="close">Close</v-btn>
          <v-btn color="blue darken-1" text @click="lock">Lock</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-row>
</template>

<script>
import { required, minLength, maxLength } from 'vuelidate/lib/validators'
export default {
  name: 'LockDialog',
  props: {
    stage: String,
    active: Boolean
  },
  data () {
    return {
      lockCode: '',
      lockOwner: '',
      lockComment: ''
    }
  },
  validations: {
    lockCode: {
      minLength: minLength(3),
      maxLength: maxLength(8)
    },
    lockOwner: {
      required,
      minLength: minLength(3),
      maxLength: maxLength(100)
    },
    lockComment: {
      maxLength: maxLength(500)
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
    async lock () {
      this.$v.$touch()
      if (this.$v.$invalid) {
        this.$emit('notify', 'Check form fields', 'red')
        return false
      }
      let data = null
      try {
        let response = await this.$http.post(`/stages/${this.stage}/lock`, {
          code: this.lockCode,
          locked_by: this.lockOwner,
          comment: this.lockComment
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
