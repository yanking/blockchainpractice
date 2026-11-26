// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.39;

// Bank 合约：接收 ETH 存款，记录每个地址的存款金额，
// 维护存款金额前 3 名，只有管理员能提取合约里的全部 ETH。
contract Bank {
    // 管理员，即部署合约的地址
    address payable public immutable admin;

    // 每个地址的存款金额
    mapping(address => uint256) public deposits;

    // 存款金额前 3 名的地址，按金额从高到低排列，不足 3 人时空位为 address(0)
    address[3] public top3;

    event Deposit(address indexed account, uint256 amount);
    event Withdraw(address indexed to, uint256 amount);

    constructor() {
        admin = payable(msg.sender);
    }

    // 存款：用 MetaMask 向合约地址转账（不带 data）就会走到这里
    receive() external payable {
        require(msg.value > 0, "Bank: deposit must be greater than 0");

        deposits[msg.sender] += msg.value;
        _updateTop3(msg.sender);
        emit Deposit(msg.sender, msg.value);
    }

    // 管理员提取合约里的全部 ETH，存款记录和前 3 名保持不变
    function withdraw() external onlyAdmin {
        uint256 amount = address(this).balance;
        (bool ok, ) = admin.call{value: amount}("");
        require(ok, "Bank: transfer failed");

        emit Withdraw(admin, amount);
    }

    modifier onlyAdmin() {
        require(msg.sender == admin, "Bank: not admin");
        _;
    }

    // 存款金额变化后更新前 3 名。
    // 存款只增不减，所以 account 在榜上时只会往前挪。
    // 金额相同时排在已入榜的地址后面。
    function _updateTop3(address account) internal {
        uint256 amount = deposits[account];
        uint256 slot = 3; // 新位置，3 表示进不了前 3
        uint256 last = 2; // 搬动上界：不在榜时最后一名被顶掉

        // 一趟扫描同时确定「是否在榜」「插到哪里」
        for (uint256 i = 0; i < 3; i++) {
            address holder = top3[i];
            if (holder == account) {
                last = i; // 在榜，只在自己上方腾位置
                break;
            }
            if (slot == 3 && (holder == address(0) || deposits[holder] < amount)) {
                slot = i;
            }
        }
        if (slot == 3) {
            return; // 已在正确位置，或进不了前 3
        }

        // 只搬 slot..last 这一段，account 原位置被前移覆盖
        for (uint256 j = last; j > slot; j--) {
            top3[j] = top3[j - 1];
        }
        top3[slot] = account;
    }
}
