//
//  PasswordManager.swift
//  Flappy Bird
//
//  Created by Brandon Plank on 9/26/21.
//  Copyright © 2021 Brandon Plank & Thatcher Clough. All rights reserved.
//

import Foundation
import FlappyEncryption

class PasswordManager {
    
    static let nameKey = "name"
    static let passwordKey = "password"
    
    public struct User {
        let name: String?
        let password: String?
    }
    
    static let defaults = UserDefaults.standard
    
    static func readUser() -> PasswordManager.User? {
        guard let encryptedName = defaults.string(forKey: nameKey) else {
            #if DEBUG
            print("Failed to read encryted username")
            #endif
            return nil
        }
        #if DEBUG
        print("Read encrypted name: \(encryptedName)")
        #endif
        
        guard let encryptedPassword = defaults.string(forKey: passwordKey) else {
            #if DEBUG
            print("Failed to read encryted password")
            #endif
            return nil
        }
        
        #if DEBUG
        print("Read encrypted password: \(encryptedPassword)")
        #endif

        
        guard let decryptedName = FlappyEncryption.decryptString(base64: encryptedName) else {
            #if DEBUG
            print("Failed to decrypt username")
            #endif
            return nil
        }
        
        guard let decryptedPassword = FlappyEncryption.decryptString(base64: encryptedPassword) else {
            #if DEBUG
            print("Failed to decrypt password")
            #endif
            return nil
        }
        
        return PasswordManager.User(name: decryptedName, password: decryptedPassword)
    }
    
    static func writeUser(_ user: PasswordManager.User) {
        guard let encryptedName = FlappyEncryption.encryptBase64String(user.name!) else {
            #if DEBUG
            print("Failed to encrypt username")
            #endif
            return
        }
        
        guard let encryptedPassword = FlappyEncryption.encryptBase64String(user.password!) else {
            #if DEBUG
            print("Failed to encrypt password")
            #endif
            return
        }
        #if DEBUG
        print("Writing encrypted username: \(encryptedName)")
        #endif
        defaults.set(encryptedName, forKey: nameKey)
        #if DEBUG
        print("Writing encrypted password: \(encryptedPassword)")
        #endif
        defaults.set(encryptedPassword, forKey: passwordKey)
        defaults.synchronize()
    }
    
    static func isSignedIn() -> Bool {
        guard let user = readUser() else {
            return false
        }
        if user.name == nil || user.password == nil {
            return false
        }
        return true
    }
}
