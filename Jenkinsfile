pipeline {
  options {
    buildDiscarder(logRotator(numToKeepStr: '7'))
  }
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
    stage('Unit Test') {
      steps {
        container('golang') {
          sh """
            mkdir -p /go/src/aheadaviation
            ln -s `pwd` /go/src/aheadaviation/Users
            cd /go/src/aheadaviation/Users
            make dep
            make test
          """
        }
      }
    }
    stage('Check Code Quality') {
      environment {
        scannerHome = tool 'SonarQube Scanner'
      }
      steps {
        withSonarQubeEnv('SonarQube') {
          sh "${scannerHome}/bin/sonar-scanner"
        }
        sleep(10)
        timeout(time: 10, unit: 'MINUTES') {
          waitForQualityGate abortPipeline: true
        }
      }
    }
    stage('Build & Push Binary') {
      steps {
        container('golang') {
          sh """
            mkdir -p /go/src/aheadaviation
            ln -s `pwd` /go/src/aheadaviation/Users
            cd /go/src/aheadaviation/Users
            make dep
            make linux
          """
        }
      }
      post {
        success {
          archiveArtifacts artifacts: 'bin/users-amd64-linux', fingerprint: true
        }
      }
    }
  }
}
