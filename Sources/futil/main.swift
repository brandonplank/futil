import ArgumentParser
import FlappyAPI
import FlappyEncryption

struct futil: ParsableCommand {
    static let configuration = CommandConfiguration(
        abstract: "Internal Flappy Bird moderation tool\nIf you do not have permission to use this, you should not."
    )
    
    @Option(name: [.customLong("name"), .customShort("n")], help: "Username")
    var name: String?
    
    @Option(name: [.customLong("password"), .customShort("p")], help: "Password")
    var password: String?
    
    @Option(name: [.customLong("score"), .customShort("s")], help: "Score")
    var score: Int?
    
    @Flag(name: [.customLong("login"), .customShort("l")], help: "Login to your admin account")
    var login = false
    
    @Flag(name: [.customLong("id"), .customShort("i")], help: "Get the UUID of a user by name, use with -n")
    var id = false
    
    @Flag(name: [.customLong("list")], help: "Lists all of the users, and their score")
    var list = false
    
    @Flag(name: [.customLong("ban"), .customShort("b")], help: "Ban a user by name, use -n")
    var ban = false
    
    @Flag(name: [.customLong("delete"), .customShort("x")], help: "Delete a user by name, use -n")
    var delete = false
    
    @Flag(name: [.customLong("unban"), .customShort("u")], help: "Unban a user by name, use -n")
    var unban = false
    
    @Flag(name: [.customLong("admin"), .customShort("a")], help: "Make a user a admin, use -n")
    var admin = false
    
    @Flag(name: [.customLong("restoreScore"), .customShort("r")], help: "Restore a users score to any value, use -n and -s")
    var restoreScore = false
    
    @Flag(name: [.customLong("count"), .customShort("c")], help: "Get the player count")
    var count = false
    
    @Flag(name: [.customLong("deaths"), .customShort("d")], help: "Get the player deaths")
    var deaths = false
    
    @Flag(name: [.customLong("verbose"), .customShort("v")], help: "Show extra logging for debugging purposes")
    var verbose = false
    
    mutating func run() throws {
        if list {
            if let user = PasswordManager.readUser() {
                let users = FlappyAPI(endpoint: "internal_users").getUsers(user.name!, user.password!)
                guard let users = users else {
                    print("Failed to get users")
                    return
                }
                for flappyuser in users {
                    print("\(flappyuser.name) | Score: \(flappyuser.score ?? 0) | Deaths: \(flappyuser.deaths!)")
                }
            } else {
                print("Login with -l")
                throw ExitCode(-1)
            }
        }
        if login {
            if verbose { print("Logging in") }
            if let name = name {
                if verbose { print("Name: \(name)") }
                if let password = password {
                    if verbose { print("Password: \(password)") }
                    PasswordManager.writeUser(PasswordManager.User(name: name, password: password))
                } else {
                    print("Use -p")
                }
            } else {
                print("Use -n")
            }
        }
        if id {
            if let name = name {
                print(FlappyAPI(endpoint: "getID").getUserID(name)!)
            }
        }
        if count {
            print("Players: \(FlappyAPI(endpoint: "userCount").getInt())")
        }
        if deaths {
            print("Deaths: \(FlappyAPI(endpoint: "globalDeaths").getInt())")
        }
        if ban {
            if let name = name {
                print("Banning \(name)")
                Utils.ban(name)
            } else {
                print("Use -n")
            }
        }
        if unban {
            if let name = name {
                print("Unbanning \(name)")
                Utils.unban(name)
            } else {
                print("Use -n")
            }
        }
        if delete {
            if let name = name {
                if let user = PasswordManager.readUser() {
                    FlappyAPI(endpoint: "delete").deleteUser(user.name!, user.password!, Utils.getUserId(name))
                    print("Deleted \(name)'s account")
                } else {
                    print("Unable to delete user")
                    throw ExitCode(-1)
                }
            }
        }
        if restoreScore {
            if let name = name {
                if let score = score {
                    if let user = PasswordManager.readUser() {
                        FlappyAPI(endpoint: "restoreScore").restoreScore(user.name!, user.password!, Utils.getUserId(name), score)
                        print("\(name)'s score restored to \(score)")
                    } else {
                        print("Unable to ban user")
                        throw ExitCode(-1)
                    }
                } else {
                    print("Use -s")
                }
            } else {
                print("Use -n")
            }
        }
        if admin {
            if let name = name {
                if let user = PasswordManager.readUser() {
                    FlappyAPI(endpoint: "makeAdmin").makeAdmin(user.name!, user.password!, Utils.getUserId(name))
                    print("Made \(name) admin")
                } else {
                    print("Unable to ban user")
                    throw ExitCode(-1)
                }
            }
        }
    }
}

futil.main()
