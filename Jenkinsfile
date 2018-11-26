pipeline {
  agent {
    kubernetes {
      label 'bagstore-users'
      defaultContainer 'jnlp'
      yaml """
apiVersion: v1
kind: Pod
metadata:
labels:
  component: ci
spec:
  containers:
  - name: golang
    image: golang:1.10
    command:
    - cat
    tty: true
"""
    }
  }
  stages {
    stage('Test') {
      steps {
        checkout scm
        container('golang') {
          sh """
            pwd
            ls
            make test
          """
        }
      }
    }
  }
}
