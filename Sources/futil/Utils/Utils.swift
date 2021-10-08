//
//  Utils.swift
//  
//
//  Created by Brandon Plank on 10/8/21.
//

import Foundation
import FlappyAPI
import FlappyEncryption

class Utils {
    
    static func getUserId(_ name: String) -> String {
        guard let id = FlappyAPI(endpoint: "getID").getUserID(name) else {
            print("Error with getting user ID")
            exit(EXIT_FAILURE)
        }
        return id
    }
    
    static func ban(_ name: String, _ reason: String? = nil) {
        if let user = PasswordManager.readUser() {
            FlappyAPI(endpoint: "ban").ban(user.name!, user.password!, getUserId(name), reason)
        } else {
            print("Unable to ban user")
            exit(EXIT_FAILURE)
        }
    }
    
    static func unban(_ name: String) {
        if let user = PasswordManager.readUser() {
            FlappyAPI(endpoint: "unban").unban(user.name!, user.password!, getUserId(name))
        } else {
            print("Unable to unban user")
            exit(EXIT_FAILURE)
        }
    }
}
